package maven

import (
	_ "github.com/lib/pq"
	"time"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"io"
)

var scripts []string = []string{
	"create table db_changelog (id bigint primary key, executedAt bigint, hash varchar(100))",
	"create sequence account_seq start 1; create table accounts (id bigint primary key default nextval('account_seq'), username varchar(255), password varchar(100), admin boolean default false)",
	"create sequence domain_seq start 1; create table domains (id bigint primary key default nextval('domain_seq'), name varchar(512), active boolean, userId bigint)",
	"create sequence package_seq start 1; create table packages (id bigint primary key default nextval('package_seq'), groupName varchar(512), password varchar(100), public boolean, domainId bigint)",
	"create sequence artifact_seq start 1; create table artifacts (id bigint primary key default nextval('artifact_seq'), groupName varchar(512), artifactName varchar(256), version varchar(128), filename varchar(128), packageId bigint not null)",
	"create unique index artifact_unq on artifacts (groupName, artifactName, version, filename)",
	"create unique index name_unq on domains (name)",
	"create unique index username_unq on accounts (username)",
	"create sequence session_seq start 1; create table sessions (id varchar(255) primary key, created bigint, expires bigint, ip varchar(255), userId bigint)",
}

func SetupDatabase() error {
	firstRun, err := firstRun()

	if err != nil {
		return err
	}

	var index int = 0

	if !firstRun {
		// select highest completed sql query.
		count, err := max()

		if err != nil {
			return err
		}

		index = count
	} else {
		// execute all scripts
		err := execute(0)

		if err != nil {
			return err
		}
	}

	if index < len(scripts) - 1 {
		// execute each script after index.
		err := execute(index + 1)

		if err != nil {
			return err
		}
	}

	return nil
}

func firstRun() (bool, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return false, err
	}

	defer conn.Close()

	var count int

	row := conn.QueryRow("select count(table_name) from information_schema.tables WHERE table_type = 'BASE TABLE' AND table_schema = 'public'")
	err = row.Scan(&count)

	if err != nil {
		return  true, err
	}

	return count == 0, nil
}

func max() (int, error) {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return 0, err
	}

	defer conn.Close()

	var count int = 0

	row := conn.QueryRow("select max(id) from db_changelog")
	err = row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func execute(from int) error {
	conn, err := sql.Open("postgres", "")

	if err != nil {
		return err
	}

	defer conn.Close()

	for index := from; index < len(scripts); index++ {
		query := scripts[index]
		_, err := conn.Exec(query)

		if err != nil {
			return err
		}

		now := time.Now().Unix()

		hash := hash(query)

		_, err = conn.Exec("insert into db_changelog (id, executedAt, hash) values ($1, $2, $3)", index, now, hash)

		if err != nil {
			return err
		}
	}

	return nil
}

func hash(sql string) string {
	hasher := sha1.New()
	io.WriteString(hasher, sql)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
