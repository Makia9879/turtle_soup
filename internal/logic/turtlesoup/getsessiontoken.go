package turtlesoup

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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
)

type GetSessionToken struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

type StoryContent struct {
	Title   string `json:"title"`
	Surface string `json:"surface"`
	Bottom  string `json:"bottom"`
}

func NewGetSessionToken(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *GetSessionToken {
	return &GetSessionToken{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *GetSessionToken) GetSessionToken(req *types.GetSessionTokenRequest) (resp *types.GetSessionTokenResponse, err error) {
	valCtx := contextx.ValueOnlyFrom(l.ctx)
	c := l.svcCtx.MustGetConfig()

	// 1. 检查活动token是否有效
	activeToken, err := l.svcCtx.RedisConn.GetCtx(valCtx, constant.GetActivityToken(req.ActivityToken))
	if err != nil {
		l.Logger.Errorf("[GetSessionToken] GetCtx() error: %v", err)
		return nil, err
	} else if activeToken == "" {
		return nil, errors.ErrActiveTokenExpired
	}

	// 2. 生成sessionToken
	sessionToken := uuid.New().String()

	// 3. 获取一个未玩过的随机海龟汤故事
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
	            AND st.activity_token = ?
	        )
	    ORDER BY 
	        RAND() 
	    LIMIT 1
	`
	err = l.svcCtx.SqlConn.QueryRowCtx(valCtx, &story, query, req.ActivityToken)
	if err != nil && pkgerrors.Is(err, t_turtle_soup_stories.ErrNotFound) {
		l.Logger.Errorf("[GetSessionToken] 没有可用的故事")
		return nil, err
	} else if err != nil {
		l.Logger.Errorf("[GetSessionToken] SqlConn.QueryRowCtx() error: %v", err)
		return nil, err
	}

	// 4. 保存sessionToken到数据库
	err = l.svcCtx.SqlConn.TransactCtx(valCtx, func(ctx context.Context, sqlSession sqlx.Session) error {
		_, err = l.svcCtx.Model.TSessionTokens.Insert(valCtx, sqlSession, &t_session_tokens.TSessionTokens{
			Token:             sessionToken,
			ActivityToken:     req.ActivityToken,
			RemainingAttempts: int64(c.DefaultRemainingTries),
			RemainingAnswers:  int64(c.DefaultRemainingAnswers),
		})
		if err != nil {
			l.Logger.Errorf("[GetSessionToken] TSessionTokens.Insert() error: %v", err)
			return err
		}

		_, err = l.svcCtx.Model.TUserSessions.Insert(valCtx, sqlSession, &t_user_sessions.TUserSessions{
			SessionToken: sessionToken,
			StoryId:      int64(story.Id),
		})
		if err != nil {
			l.Logger.Errorf("[GetSessionToken] TUserSessions.Insert() error: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		l.Logger.Errorf("[GetSessionToken] SqlConn.TransactCtx() error: %v", err)
		return nil, err
	}

	// 5. 保存到缓存
	storyContent := StoryContent{
		Title:   story.Title,
		Surface: story.Surface,
		Bottom:  story.Bottom,
	}
	jsonContent, err := json.Marshal(storyContent)
	sessionTokenCacheKey := constant.GetSessionToken(sessionToken)
	err = l.svcCtx.RedisConn.HmsetCtx(valCtx, sessionTokenCacheKey, map[string]string{
		"remainingTries":   strconv.Itoa(c.DefaultRemainingTries),
		"remainingAnswers": strconv.Itoa(c.DefaultRemainingAnswers),
		"storyId":          strconv.FormatUint(story.Id, 10),
		"storyContent":     string(jsonContent),
	})
	if err != nil {
		l.Logger.Errorf("[GetSessionToken] HmsetCtx() error: %v", err)
		return nil, err
	}

	// 6. 返回响应
	return &types.GetSessionTokenResponse{
		Token:            sessionToken,
		RemainingTries:   c.DefaultRemainingTries,
		RemainingAnswers: c.DefaultRemainingAnswers,
		StoryTitle:       story.Title,
		StoryID:          int(story.Id),
		Surface:          story.Surface,
	}, nil
}
