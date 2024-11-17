package db

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate ../../bin/minimock -i Transactor -o ./mocks/ -s "_minimock.go"
//go:generate ../../bin/minimock -i TxManager -o ./mocks/ -s "_minimock.go"
//go:generate ../../bin/minimock -i github.com/jackc/pgx/v4.Tx -o ./mocks/ -s "_minimock.go"
