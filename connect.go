package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"bitbucket.org/kardianos/table"
	_ "github.com/lib/pq"
)

func main() {
	args := strings.Split("start --insecure --listen-addr=localhost:9999 --store=type=mem,size=1GiB --logtostderr=NONE", " ")
	cmd := exec.Command("cockroach", args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer cmd.Process.Signal(os.Interrupt)

	db, err := sql.Open("postgres", "postgresql://root@localhost:9999?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	for i := 0; i < 3; i++ {
		err = db.PingContext(ctx)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
	if err != nil {
		log.Fatal("failed to ping db", err)
	}

	set, err := table.NewSet(ctx, db, `
		begin;
			select * from information_schema.tables limit 10;
			select * from information_schema.columns limit 10;
		commit;
	`)
	if err != nil {
		log.Fatal("failed to fill set", err)
	}
	for tindex, t := range set {
		fmt.Printf("\nset %d\n", tindex+1)
		for ci, c := range t.Columns {
			if ci != 0 {
				fmt.Print(",")
			}
			fmt.Printf("%s", c)
		}
		fmt.Println()
		for _, row := range t.Rows {
			for ci, c := range t.Columns {
				if ci != 0 {
					fmt.Print(",")
				}
				fmt.Printf("%v", row.Value(c))
			}
			fmt.Println()
		}
	}

	log.Println("done (waiting for shutdown...)")
	cmd.Wait()
	log.Println("end")
}
