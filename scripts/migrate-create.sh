if [ -z "$1" ]
then
    echo "need migration name"
else
    migrate create -ext sql -dir ./migrations -seq $1
fi
