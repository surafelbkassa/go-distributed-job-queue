-- KEYS[1]=queue:delayed KEYS[2]=queue:pending
-- ARGV[1]=limit
local t = redis.call('TIME')
local ids = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', tonumber(t[1]), 'LIMIT', 0, ARGV[1])
local n = 0
for _, id in ipairs(ids) do
  if redis.call('ZREM', KEYS[1], id) == 1 then
    redis.call('LPUSH', KEYS[2], id)
    redis.call('HSET', 'jobs:' .. id, 'status', 'pending')
    n = n + 1
  end
end
return n