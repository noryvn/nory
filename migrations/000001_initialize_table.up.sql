BEGIN;

CREATE TABLE IF NOT EXISTS app_user (
	user_id UUID UNIQUE NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	username VARCHAR(20) NOT NULL,
	name VARCHAR(32) NOT NULL,
	email VARCHAR(254) UNIQUE NOT NULL,

	PRIMARY KEY(user_id)
);

-- some times username will be used in WHERE clause
CREATE INDEX IF NOT EXISTS user_id_index ON app_user(user_id, username);
CREATE UNIQUE INDEX lower_username ON app_user(LOWER(username));

CREATE TABLE IF NOT EXISTS class (
	class_id VARCHAR(20) UNIQUE NOT NULL,
	owner_id UUID NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	name VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,

	PRIMARY KEY(class_id),
	CONSTRAINT fk_user FOREIGN KEY (owner_id) REFERENCES app_user(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_index ON class(class_id);

CREATE TABLE IF NOT EXISTS class_schedule (
	schedule_id VARCHAR(20) UNIQUE NOT NULL,
	class_id VARCHAR(20) NOT NULL,

	name VARCHAR(20) NOT NULL,
	start_at TIME NOT NULL,
	duration SMALLINT NOT NULL,
	day SMALLINT NOT NULL,

	PRIMARY KEY(schedule_id),
	CONSTRAINT fk_class FOREIGN KEY (class_id) REFERENCES class(class_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_schedule_index ON class_schedule(class_id);

CREATE TABLE IF NOT EXISTS class_task (
	task_id VARCHAR(20) UNIQUE NOT NULL,
	class_id VARCHAR(20) NOT NULL,

	name VARCHAR(20) NOT NULL,
	due_date TIMESTAMP NOT NULL,
	description VARCHAR(1024) NOT NULL,

	PRIMARY KEY(task_id),
	CONSTRAINT fk_class FOREIGN KEY (class_id) REFERENCES class(class_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_task_index ON class_task(class_id);

COMMIT;
