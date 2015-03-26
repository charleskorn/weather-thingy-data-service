-- +migrate Up
CREATE TABLE variables (
  variable_id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  units VARCHAR(20) NOT NULL,
  created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE variables;
