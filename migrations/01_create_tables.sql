CREATE TABLE project (
    jira_id BIGINT PRIMARY KEY,
    key VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    url TEXT
);

CREATE TABLE author (
    jira_id BIGINT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255)
);

CREATE TABLE issue (
    jira_id BIGINT PRIMARY KEY,
    project_id BIGINT NOT NULL,
    key VARCHAR(50) NOT NULL,
    summary TEXT,
    status VARCHAR(100),
    priority VARCHAR(50),
    created_time TIMESTAMP NOT NULL,
    updated_time TIMESTAMP,
    closed_time TIMESTAMP NULL,
    time_spent INTEGER,
    creator_id BIGINT,
    assignee_id BIGINT,
    FOREIGN KEY (project_id) REFERENCES project(jira_id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES author(jira_id),
    FOREIGN KEY (assignee_id) REFERENCES author(jira_id),
    CONSTRAINT unique_issue_key UNIQUE(project_id, key)
);

CREATE TABLE status_change (
    id BIGSERIAL PRIMARY KEY,
    issue_id BIGINT NOT NULL,
    old_status VARCHAR(100),
    new_status VARCHAR(100),
    change_time TIMESTAMP NOT NULL,
    FOREIGN KEY (issue_id) REFERENCES issue(jira_id) ON DELETE CASCADE,
    CONSTRAINT unique_status_change UNIQUE (issue_id, change_time, new_status)
);

CREATE INDEX idx_project_key ON project(key);
CREATE INDEX idx_author_username ON author(username);
CREATE INDEX idx_issues_project_id ON issue(project_id);
CREATE INDEX idx_issues_status ON issue(status);
CREATE INDEX idx_issues_priority ON issue(priority);
CREATE INDEX idx_issues_created_time ON issue(created_time);
CREATE INDEX idx_issues_closed_time ON issue(closed_time);
CREATE INDEX idx_statuschange_issue_id ON status_change(issue_id);
CREATE INDEX idx_statuschange_change_time ON status_change(change_time);