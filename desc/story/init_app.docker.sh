
# 遍历/desc/story/下所有json文件并执行添加命令
for file in /dist/desc/story/*.json; do
    /dist/app story add -f "$file"
done
