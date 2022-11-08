BEGIN;

CREATE TABLE IF NOT EXISTS app_user (
	user_id UUID UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	username VARCHAR(20) NOT NULL,
	name VARCHAR(32) NOT NULL,
	email VARCHAR(254) NOT NULL,

	CONSTRAINT app_user_pk PRIMARY KEY(user_id)
);

CREATE INDEX IF NOT EXISTS user_id_index ON app_user(user_id);
CREATE INDEX IF NOT EXISTS username_index ON app_user(username);
CREATE UNIQUE INDEX lower_username ON app_user(LOWER(username));
CREATE UNIQUE INDEX lower_email ON app_user(LOWER(email));

CREATE TABLE IF NOT EXISTS class (
	class_id VARCHAR(20) UNIQUE NOT NULL,
	owner_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	name VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,

	CONSTRAINT class_pk PRIMARY KEY(class_id),
	CONSTRAINT fk_user FOREIGN KEY (owner_id) REFERENCES app_user(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_id_index ON class(class_id);
CREATE INDEX IF NOT EXISTS owner_id_index ON class(owner_id, class_id);
CREATE UNIQUE INDEX user_class_name ON class (owner_id, name);

CREATE TABLE IF NOT EXISTS class_schedule (
	schedule_id VARCHAR(20) UNIQUE NOT NULL,
	class_id VARCHAR(20) NOT NULL,
	author_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	name VARCHAR(20) NOT NULL,
	start_at TIME NOT NULL,
	duration SMALLINT NOT NULL,
	day SMALLINT NOT NULL,

	CONSTRAINT class_schedule_pk PRIMARY KEY(schedule_id),
	CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES app_user(user_id) ON DELETE CASCADE,
	CONSTRAINT fk_class FOREIGN KEY (class_id) REFERENCES class(class_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_schedule_index ON class_schedule(schedule_id);
CREATE INDEX IF NOT EXISTS class_id_index ON class_schedule(class_id, schedule_id);

CREATE TABLE IF NOT EXISTS class_task (
	task_id VARCHAR(20) UNIQUE NOT NULL,
	class_id VARCHAR(20) NOT NULL,
	author_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	author_display_name VARCHAR(20) NOT NULL,
	name VARCHAR(20) NOT NULL,
	due_date DATE NOT NULL,
	description VARCHAR(1024) NOT NULL,

	CONSTRAINT class_task_pk PRIMARY KEY(task_id),
	CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES app_user(user_id) ON DELETE CASCADE,
	CONSTRAINT fk_class FOREIGN KEY (class_id) REFERENCES class(class_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_task_index ON class_task(task_id, due_date);
CREATE INDEX IF NOT EXISTS class_id_index ON class_task(class_id, task_id);

CREATE TABLE IF NOT EXISTS class_member (
	class_id VARCHAR(20) NOT NULL,
	user_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	level VARCHAR(10) NOT NULL,

	CONSTRAINT class_member_pk PRIMARY KEY(class_id, user_id),
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES app_user(user_id) ON DELETE CASCADE,
	CONSTRAINT fk_class FOREIGN KEY (class_id) REFERENCES class(class_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_id_index ON class_member(class_id, created_at);
CREATE INDEX IF NOT EXISTS user_id_index ON class_member(user_id, created_at);
CREATE INDEX IF NOT EXISTS user_and_class_index ON class_member(user_id, class_id, created_at);

COMMIT;
