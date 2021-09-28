CREATE TABLE cb_user_states (
	user_id INT NOT NULL,
	related_to TIMESTAMP NOT NULL,
	level INT NOT NULL,
	last_update TIMESTAMP NOT NULL,

	ancient_shard INT  NOT NULL,
	void_shard INT  NOT NULL,
	sacred_shard INT  NOT NULL,
	epic_tome INT  NOT NULL,
	leg_tome INT  NOT NULL,
	PRIMARY KEY(user_id, related_to, level)
);

CREATE INDEX cb_user_states_lu_idx ON cb_user_states(last_update);
CREATE INDEX cb_user_states_rel_idx ON cb_user_states(related_to);
CREATE INDEX cb_user_user_id_idx ON cb_user_states(user_id);