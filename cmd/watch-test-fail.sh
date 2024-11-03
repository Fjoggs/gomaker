while true
do
    inotifywait -qq -r -e create,close_write,modify,move,delete ./ && grc go test -v ./... | grep FAIL -B 10
done
