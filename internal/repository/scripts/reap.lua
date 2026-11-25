-- KEYS[1]=queue:processing
-- ARGV[1]=limit
local t = redis.call('TIME')
return redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', tonumber(t[1]), 'LIMIT', 0, ARGV[1])