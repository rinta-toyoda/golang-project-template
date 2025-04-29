CREATE TABLE user (
    id CHAR(36),
    email VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE user_profile (
    id CHAR(36),
    user_id CHAR(36),
    first_name VARCHAR(20) NOT NULL,
    middle_name VARCHAR(20),
    last_name VARCHAR(20) NOT NULL,
    address VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user(id)
);

CREATE DOMAIN years AS NUMERIC(3, 1) CHECK(VALUE >= 0);

CREATE TABLE experience (
    id CHAR(36),
    profile_id CHAR(36),
    company VARCHAR(20) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(id),
    FOREIGN KEY(profile_id) REFERENCES user_profile(id)
);

CREATE TABLE experience_skill (
    id CHAR(36),
    experience_id CHAR(36),
    name VARCHAR(20) NOT NULL,
    years years NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(experience_id) REFERENCES experience(id),
    UNIQUE(experience_id, name)
);

CREATE TABLE education (
    id CHAR(36),
    profile_id CHAR(36),
    institution VARCHAR(20) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(id),
    FOREIGN KEY(profile_id) REFERENCES user_profile(id)
);

CREATE TABLE education_skill (
    id CHAR(36),
    education_id CHAR(36),
    name VARCHAR(20) NOT NULL,
    years years NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(education_id) REFERENCES education(id),
    UNIQUE(education_id, name)
);

CREATE TABLE organization (
    id CHAR(36),
    email VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(15) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE organization_profile (
    id CHAR(36),
    organization_id CHAR(36),
    name VARCHAR(20) NOT NULL,
    description TEXT,
    address VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(id),
    FOREIGN KEY(organization_id) REFERENCES organization(id)
);

CREATE TABLE job (
    id CHAR(36),
    profile_id CHAR(36),
    name VARCHAR(20) NOT NULL,
    description TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(id),
    FOREIGN KEY(profile_id) REFERENCES organization_profile(id)
);

CREATE TABLE required_skill (
    id CHAR(36),
    job_id CHAR(36),
    name VARCHAR(20) NOT NULL,
    years years NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(job_id) REFERENCES job(id)
);

CREATE TABLE application (
    id CHAR(36),
    user_id CHAR(36),
    job_id CHAR(36),
    applied_at TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES user(id),
    FOREIGN KEY(job_id) REFERENCES job(id),
    UNIQUE(user_id, job_id)
);

CREATE TABLE application_skill (
    id CHAR(36),
    application_id CHAR(36),
    name VARCHAR(50) NOT NULL,
    years years NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(id),
    FOREIGN KEY(application_id) REFERENCES application(id)
);


CREATE INDEX idx_user_email ON user(email);
CREATE INDEX idx_application_user_job ON application(user_id, job_id);
CREATE INDEX idx_userprofile_user_id ON user_profile(user_id);
