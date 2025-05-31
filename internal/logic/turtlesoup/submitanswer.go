package turtlesoup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/jzero-io/jzero-contrib/condition"
	pkgerrors "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"turtle-soup/internal/constant"
	"turtle-soup/internal/errors"
	"turtle-soup/internal/model/t_session_tokens"
	"turtle-soup/internal/model/t_turtle_soup_stories"
	"turtle-soup/internal/model/t_user_sessions"
	"turtle-soup/internal/svc"
	types "turtle-soup/internal/types/turtlesoup"
	"turtle-soup/pkg/deepseek"
)

type SubmitAnswer struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewSubmitAnswer(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *SubmitAnswer {
	return &SubmitAnswer{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *SubmitAnswer) handleNewStoryRequest(valCtx context.Context, sessionToken string, sessionTokenCacheKey string, remainingTries int, remainingAnswers int) (*types.SubmitAnswerResponse, error) {
	// 获取一个未玩过的随机海龟汤故事
	var story t_turtle_soup_stories.TTurtleSoupStories
	query := `
	    SELECT 
	        ts.*
	    FROM 
	        T_turtle_soup_stories ts 
	    WHERE 
	        ts.id NOT IN (
	            SELECT story_id 
	            FROM T_user_sessions us 
	            LEFT JOIN T_session_tokens st ON st.token = us.session_token 
	            WHERE is_completed = FALSE 
	            AND st.activity_token = (
	                SELECT activity_token FROM T_session_tokens WHERE token = ?
	            )
	        )
	    ORDER BY 
	        RAND() 
	    LIMIT 1
	`
	err := l.svcCtx.SqlConn.QueryRowCtx(valCtx, &story, query, sessionToken)
	if err != nil && pkgerrors.Is(err, t_turtle_soup_stories.ErrNotFound) {
		l.Logger.Errorf("[handleNewStoryRequest] 没有可用的故事")
		return nil, err
	} else if err != nil {
		l.Logger.Errorf("[handleNewStoryRequest] SqlConn.QueryRowCtx() error: %v", err)
		return nil, err
	}

	// 更新用户会话和缓存
	err = l.svcCtx.SqlConn.TransactCtx(valCtx, func(ctx context.Context, sqlSession sqlx.Session) error {
		_, err = l.svcCtx.Model.TUserSessions.Insert(valCtx, sqlSession, &t_user_sessions.TUserSessions{
			SessionToken: sessionToken,
			StoryId:      int64(story.Id),
		})
		if err != nil {
			l.Logger.Errorf("[handleNewStoryRequest] TUserSessions.Update() error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		l.Logger.Errorf("[handleNewStoryRequest] SqlConn.TransactCtx() error: %v", err)
		return nil, err
	}

	// 更新缓存
	storyContent := StoryContent{
		Title:   story.Title,
		Surface: story.Surface,
		Bottom:  story.Bottom,
	}
	jsonContent, err := json.Marshal(storyContent)
	err = l.svcCtx.RedisConn.HmsetCtx(valCtx, sessionTokenCacheKey, map[string]string{
		"remainingTries":   strconv.Itoa(remainingTries),
		"remainingAnswers": strconv.Itoa(remainingAnswers),
		"storyId":          strconv.FormatUint(story.Id, 10),
		"storyContent":     string(jsonContent),
	})
	if err != nil {
		l.Logger.Errorf("[handleNewStoryRequest] HmsetCtx() error: %v", err)
		return nil, err
	}

	return &types.SubmitAnswerResponse{
		Reply:            "开启新的故事",
		RemainingTries:   remainingTries,
		RemainingAnswers: remainingAnswers,
		NewStoryId:       int64(story.Id),
		NewStoryTitle:    story.Title,
		NewSurface:       story.Surface,
	}, nil
}

func (l *SubmitAnswer) SubmitAnswer(req *types.SubmitAnswerRequest) (resp *types.SubmitAnswerResponse, err error) {
	sessionToken := req.SessionToken
	messages := req.Messages
	valCtx := contextx.ValueOnlyFrom(l.ctx)
	c := l.svcCtx.MustGetConfig()

	// 检查 sessionToken绑定的 activeToken是否过期
	sessionTokenInfo, err := l.svcCtx.Model.TSessionTokens.FindOneByCondition(valCtx, nil,
		condition.NewChain().Equal("token", sessionToken).Build()...)
	if err != nil && pkgerrors.Is(err, t_session_tokens.ErrNotFound) {
		return nil, errors.ErrSessionTokenNotFound
	} else if err != nil {
		l.Logger.Errorf("查找 session token 记录时发生错误: %v", err)
		return nil, err
	}
	activeToken := sessionTokenInfo.ActivityToken
	activeTokenCacheVal, err := l.svcCtx.RedisConn.GetCtx(valCtx, constant.GetActivityToken(activeToken))
	if err != nil {
		l.Logger.Errorf("[GetSessionToken] GetCtx() error: %v", err)
		return nil, err
	} else if activeTokenCacheVal == "" {
		return nil, errors.ErrActiveTokenExpired
	}

	// 1. 检查redis中的回答次数和尝试次数
	sessionTokenCacheKey := constant.GetSessionToken(sessionToken)
	redisData, err := l.svcCtx.RedisConn.HgetallCtx(valCtx, sessionTokenCacheKey)
	if err != nil {
		l.Logger.Errorf("[SubmitAnswer] HgetallCtx() error(%s): %v", sessionToken, err)
		return nil, err
	} else if len(redisData) == 0 {
		l.Logger.Errorf("[SubmitAnswer] HgetallCtx() empty redisData(%s)", sessionToken)
		return nil, errors.ErrSessionTokenNotFound
	}

	remainingAnswers, _ := strconv.Atoi(redisData["remainingAnswers"])
	remainingTries, _ := strconv.Atoi(redisData["remainingTries"])

	// 2. 检查并更新回答次数和尝试次数
	needNewStory := false
	if remainingAnswers == 0 {
		if remainingTries == 0 {
			// 删除session_token缓存
			defer func() {
				err := l.deleteSessionToken(sessionToken)
				if err != nil {
					logx.WithContext(contextx.ValueOnlyFrom(l.ctx)).Errorf("删除 sessionToken 缓存失败: %v", err)
					return
				}
			}()
			return nil, errors.ErrNoMoreAttempts
		}
		remainingAnswers = c.DefaultRemainingAnswers
		remainingTries--
		needNewStory = true
	} else {
		remainingAnswers--
	}
	// 如果需要新的故事
	if needNewStory {
		return l.handleNewStoryRequest(valCtx, sessionToken, sessionTokenCacheKey, remainingTries, remainingAnswers)
	}

	// TODO 更新sessionToken数据库中剩余次数
	// 3. 更新redis中的次数
	err = l.svcCtx.RedisConn.HmsetCtx(valCtx, sessionTokenCacheKey, map[string]string{
		"remainingAnswers": strconv.Itoa(remainingAnswers),
		"remainingTries":   strconv.Itoa(remainingTries),
	})
	if err != nil {
		l.Logger.Errorf("[SubmitAnswer] HmsetCtx() error: %v", err)
		return nil, err
	}

	// 4. 获取海龟汤故事并构造系统消息
	storyContent := redisData["storyContent"]
	var storyContentObj StoryContent
	err = json.Unmarshal([]byte(storyContent), &storyContentObj)
	if err != nil {
		l.Logger.Errorf("[SubmitAnswer] json.Unmarshal() error: %v", err)
		return nil, err
	}
	systemMsg := types.SubmitAnswerMessage{
		Role:    "user",
		Content: c.SystemMessageTpl + storyContent + "\n请你立刻以主持人的身份开始游戏，并只展示“汤面”，不显示汤底，等待玩家开始提问。",
	}
	firstAssMsg := types.SubmitAnswerMessage{
		Role: "assistant",
		Content: fmt.Sprintf(`（主持人模式启动）
汤面：
%s
——
请开始提问、猜测或直接公布答案！
（记住：我只能回答“是”“不是”或“不重要”哦）`, storyContentObj.Surface),
	}
	lastMsg := messages[len(messages)-1].Content
	lastMsg = lastMsg + "\n以上，我的推理如果完全正确，请恭喜我完成游戏"
	messages[len(messages)-1].Content = lastMsg
	messages = append([]types.SubmitAnswerMessage{systemMsg, firstAssMsg}, messages...)

	// 5. 调用deepseek API处理消息
	deepseekMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		deepseekMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
		l.Logger.Infof("deepseekMessages[%d]%+v", i, deepseekMessages[i])
	}
	deepseekResp, err := l.svcCtx.DeepSeekClient.ChatCompletion(valCtx, &deepseek.ChatCompletionRequest{
		Messages:    deepseekMessages,
		Model:       c.DeepSeekModel,
		MaxTokens:   c.DeepSeekMaxTokens,
		Temperature: c.DeepSeekTemperature,
	})
	if err != nil {
		l.Logger.Errorf("[SubmitAnswer] DeepSeekClient.ChatCompletion() error: %v", err)
		return nil, err
	}

	// 6. 返回响应
	replyStr := deepseekResp.Choices[0].Message.Content
	isCorrect := strings.Contains(replyStr, "恭喜")
	l.Logger.Infof("replyStr: %v", replyStr)
	if !isCorrect {
		return &types.SubmitAnswerResponse{
			Reply:            replyStr,
			RemainingTries:   remainingTries,
			RemainingAnswers: remainingAnswers,
			IsCorrect:        isCorrect,
		}, nil
	}

	// 回答正确的话
	curStoryId, _ := strconv.Atoi(redisData["storyId"])
	l.svcCtx.Model.TUserSessions.UpdateFieldsByCondition(valCtx, nil,
		map[string]any{
			"is_completed": 1,
		},
		condition.NewChain().
			Equal("session_token", sessionToken).
			Equal("story_id", curStoryId).Build()...)

	// 删除session_token缓存
	defer func() {
		err := l.deleteSessionToken(sessionToken)
		if err != nil {
			logx.WithContext(contextx.ValueOnlyFrom(l.ctx)).Errorf("删除 sessionToken 缓存失败: %v", err)
			return
		}
	}()
	return &types.SubmitAnswerResponse{
		Reply:            replyStr,
		RemainingTries:   remainingTries,
		RemainingAnswers: remainingAnswers,
		IsCorrect:        true,
		StoryAnswer:      storyContentObj.Bottom,
	}, nil
}

func (l *SubmitAnswer) deleteSessionToken(sessionToken string) error {
	_, err := l.svcCtx.RedisConn.DelCtx(contextx.ValueOnlyFrom(l.ctx), constant.GetSessionToken(sessionToken))
	return err
}
