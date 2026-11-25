-- KEYS[1]=queue:pending  KEYS[2]=queue:processing
-- ARGV[1]=lease seconds
local id = redis.call('RPOP', KEYS[1])
if not id then return false end
local t = redis.call('TIME')
local deadline = tonumber(t[1]) + tonumber(ARGV[1])
redis.call('ZADD', KEYS[2], deadline, id)
local key = 'jobs:' .. id
redis.call('HINCRBY', key, 'attempts', 1)
redis.call('HSET', key, 'status', 'processing', 'updated_at', t[1])
return id