CREATE TABLE IF NOT EXISTS user (
	user_id VARCHAR(20) PRIMARY KEY,
	created_at TIMESTAMP DEFAULT NOW(),

	username VARCHAR(20) UNIQUE NOT NULL,
	name VARCHAR(20) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(64) NOT NULL
);

-- some times username will be used in WHERE clause
CREATE INDEX IF NOT EXISTS user_id_index ON user(user_id, username);

CREATE TABLE IF NOT EXISTS class (
	class_id VARCHAR(20) PRIMARY KEY,
	owner_id VARCHAR(20) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),

	name VARCHAR(20) NOT NULL,
	description VARCHAR(255) NOT NULL,

	CONTRAINT fk_user FOREIGN KEY(owner_id) REFERENCES user(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS class_index ON class(class_id);

CREATE TABLE IF NOT EXISTS class_schedule (
	schedule_id VARCHAR(20) PRIMARY KEY,
	class_id VARCHAR(20) NOT NULL,

	name VARCHAR(20) NOT NULL,
	start_at VARCHAR(5) NOT NULL,
	duration SMALLINT NOT NULL,
	day SMALLINT NOT NULL,

	CONTRAINT fk_class FOREIGN KEY(class_id) REFERENCES class(class_id) ON DELETE CASCADE
);
 ru
CREATE INDEX IF NOT EXISTS class_schedule_index ON class_schedule(class_id);
