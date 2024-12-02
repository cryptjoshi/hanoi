# wait-for-mysql.sh
until nc -z -w5 $DB_HOST $DB_PORT; do
  echo "Waiting for MySQL to be ready..."
  sleep 1
done

echo "MySQL is ready!"