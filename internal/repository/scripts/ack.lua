-- KEYS[1]=queue:processing
-- ARGV[1]=id
if redis.call('ZREM', KEYS[1], ARGV[1]) == 1 then
  redis.call('HSET', 'jobs:' .. ARGV[1], 'status', 'done')
  return 1
end
return 0