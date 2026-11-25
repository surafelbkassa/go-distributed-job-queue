-- KEYS[1]=queue:processing KEYS[2]=queue:delayed KEYS[3]=queue:dead
-- ARGV[1]=id ARGV[2]=error message ARGV[3]=backoff seconds
if redis.call('ZREM', KEYS[1], ARGV[1]) == 0 then return 'stale' end
local key = 'jobs:' .. ARGV[1]
local attempts = tonumber(redis.call('HGET', key, 'attempts')) or 0
local maxAttempts = tonumber(redis.call('HGET', key, 'max_attempts')) or 3
redis.call('HSET', key, 'last_error', ARGV[2])
if attempts >= maxAttempts then
  redis.call('LPUSH', KEYS[3], ARGV[1])
  redis.call('HSET', key, 'status', 'dead')
  return 'dead'
end
local t = redis.call('TIME')
redis.call('ZADD', KEYS[2], tonumber(t[1]) + tonumber(ARGV[3]), ARGV[1])
redis.call('HSET', key, 'status', 'retrying')
return 'retry'