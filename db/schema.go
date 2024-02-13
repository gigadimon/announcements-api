package db

var schema string = `
CREATE TABLE IF NOT EXISTS users (
    id       serial not null unique,
	login    varchar(255) not null unique,
	email    varchar(255) not null unique,
	password varchar(255) not null
);

CREATE TABLE IF NOT EXISTS announcements (
   	id           serial not null unique,
	author_id    int references users (id) on delete cascade not null,
	author_login varchar(255) references users (login) on delete cascade not null,
	author_email varchar(255) references users (email) on delete cascade not null,
	author_phone varchar(255), 
	title        varchar(255) not null,
	description  text not null,
	created_at   timestamp not null
);
`
