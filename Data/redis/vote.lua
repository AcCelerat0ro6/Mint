local personalKey = KEYS[1]
local scoreKey = KEYS[2]
local userID = ARGV[1]
local curValue = tonumber(ARGV[2])
local postID = ARGV[3]

-- 1. 获取过去的值
local pastValue = redis.call('ZSCORE', personalKey, userID)
if not pastValue then
    pastValue = 0
else
    pastValue = tonumber(pastValue)
end

-- 2. 比较过去的值和现在的值
if pastValue == curValue then
    return 0 -- 表示没做任何修改
end

-- 3. 判断是否是重复点赞/踩 (防止 1 -> 1 或者 -1 -> -1 这种非法跳跃，虽然前置相等判断拦截了，但为了更严谨)
-- 正常投票逻辑：
-- 取消投票 (1 -> 0 或 -1 -> 0)
-- 重新投票 (0 -> 1 或 0 -> -1)
-- 反转投票 (1 -> -1 或 -1 -> 1)
-- 严谨起见，如果前后都在投同一种票（比如从没投票0直接跳到投两票2，或者超出 -1~1 的范围），这应该在业务层被拦截。
-- 这里的 Lua 脚本默认接受到的 curValue 只有 -1, 0, 1，且与 pastValue 不同。

-- 4. 计算分差 (假设每一票相差 432 分)
local diff = (curValue - pastValue) * 432

-- 5. 更新帖子总分
redis.call('ZINCRBY', scoreKey, diff, postID)

-- 6. 更新个人记录
if curValue == 0 then
    redis.call('ZREM', personalKey, userID)
else
    redis.call('ZADD', personalKey, curValue, userID)
end

return 1 -- 表示修改成功