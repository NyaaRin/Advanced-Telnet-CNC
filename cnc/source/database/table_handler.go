package database

func CreateUserTable() error {
	_, err := Instance.Exec(`
			create table if not exists users
			(
			    id       INTEGER
			        constraint users_pk
			            primary key autoincrement,
			    username TEXT,
			    password TEXT,
			    methods TEXT,
			    cooldown INTEGER,
			    max_time INTEGER,
			    max_attacks INTEGER,
			    devices INTEGER,
			    admin INTEGER,
			    reseller INTEGER,
			    expiry INTEGER
			);
    `)
	return err
}

func CreateLogsTable() error {
	_, err := Instance.Exec(`
			create table if not exists logs
			(
			    id       INTEGER
			        constraint logs_pk
			            primary key autoincrement,
			    user_id INTEGER,
			    target TEXT,
			    duration INTEGER,
			    method TEXT,
			    time_created INTEGER,
			    time_end INTEGER
			);
    `)
	return err
}
