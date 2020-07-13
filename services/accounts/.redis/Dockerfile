FROM redis
COPY ./services/accounts/.redis/.redis.conf /usr/local/etc/redis/redis.conf
USER redis
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]