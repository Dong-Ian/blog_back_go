package queries

var CreateUserTable = `
	CREATE TABLE IF NOT EXISTS user_table (
		user_id 		VARCHAR(50)			NOT NULL	PRIMARY KEY,
		user_email 		VARCHAR(100)		NOT NULL	UNIQUE,
		user_password 	VARCHAR(200)		NOT NULL,
		user_name 		VARCHAR(50)		NOT NULL,
		user_status		TINYINT(1)			NOT NULL 	DEFAULT 1 COMMENT '0: 비활성, 1: 활성, 2: 탈퇴',
		profile_seq		BIGINT(20)			NULL,
		background_seq  BIGINT(20)			NULL,
		preferred_color	VARCHAR(10)			NOT NULL DEFAULT "#000",
		title			VARCHAR(50)			NOT NULL DEFAULT "Archive",
		sns_instagram	VARCHAR(100)		NULL,
		github_url		VARCHAR(100)		NULL,
		personal_url	VARCHAR(100)		NULL,
		memo			VARCHAR(200)		NULL,
		reg_date 		DATETIME			NOT NULL	DEFAULT CURRENT_TIMESTAMP,
		mod_date		DATETIME		    NULL 		DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
`

var CreatePostTable = `
	CREATE TABLE IF NOT EXISTS post_table (
		post_seq 		INT(20)		    NOT NULL AUTO_INCREMENT PRIMARY KEY,
		user_id 		VARCHAR(50)		NOT NULL REFERENCES user_table(user_id),
		post_title 		VARCHAR(50)		NOT NULL,
		post_contents 	TEXT			NOT NULL,
		post_status		TINYINT(1)		NOT NULL DEFAULT 1	COMMENT '0: 비활성, 1: 활성, 2: 삭제',
		tags			VARCHAR(150)	NULL,
		viewed			INT(20)			NOT NULL DEFAULT 0,
		is_pinned		TINYINT(1)		NOT NULL DEFAULT 1 COMMENT '0 - 비고정, 1 - 고정',
		reg_date 		DATETIME		NOT NULL DEFAULT CURRENT_TIMESTAMP,
		mod_date		DATETIME	    NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		INDEX post_idx(post_status)
	);
`

var CreateFileTable = `
	CREATE TABLE IF NOT EXISTS file_table (
		file_seq 		BIGINT(20) 			NOT NULL 	AUTO_INCREMENT PRIMARY KEY,
		file_type		TINYINT(1)			NOT NULL 	COMMENT '1 - 이미지, 2 - 비디오',
		file_format		VARCHAR(20)			NOT NULL	COMMENT 'image/png, video/avi...',
		target_seq 		INT(20)				NULL	 	COMMENT '이미지 타겟 포스트 시퀀스' REFERENCES post_table(post_seq),
		target_id		VARCHAR(50)			NOT NULL	COMMENT '이미지 타겟 유저 id - 프로필 / 업로드 유저' REFERENCES user_table(user_id),
		target_table 	VARCHAR(20)			NOT NULL,
		file_size		INT(20)				NOT NULL,
		object_name 	VARCHAR(200)		NOT NULL,
		reg_date 		DATETIME			NOT NULL	DEFAULT CURRENT_TIMESTAMP,

		INDEX target_idx(target_table, target_id)
	);
`

var CreateCommentTable = `
		CREATE TABLE IF NOT EXISTS comment_table (
			comment_seq		INT(20)			NOT NULL AUTO_INCREMENT PRIMARY KEY,
			post_seq		INT(20)			NOT NULL REFERENCES post_table(post_seq),
			user_id			VARCHAR(20)		NOT NULL REFERENCES user_table(user_id),
			COMMENT			TEXT			NOT NULL,
			reg_date 		DATETIME		NOT NULL DEFAULT CURRENT_TIMESTAMP,
			mod_date		DATETIME	    NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		)
`

//	mod_date    DATETIME        NULL        DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
