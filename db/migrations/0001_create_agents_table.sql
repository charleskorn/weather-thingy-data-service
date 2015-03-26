-- +migrate Up
CREATE TABLE agents (
  agent_id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE agents;
