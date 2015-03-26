-- +migrate Up
CREATE TABLE data (
  agent_id INT NOT NULL REFERENCES agents (agent_id),
  variable_id INT NOT NULL REFERENCES variables (variable_id),
  time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  value NUMERIC(10, 4) NOT NULL,
  PRIMARY KEY (agent_id, variable_id, time)
);

-- +migrate Down
DROP TABLE data;
