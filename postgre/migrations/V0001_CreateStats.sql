CREATE TABLE user_cb_stats (
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

CREATE INDEX user_cb_stat_lu_idx ON user_cb_stats(last_update);
CREATE INDEX user_cb_stat_rel_idx ON user_cb_stats(related_to);
CREATE INDEX user_cb_stat_id_idx ON user_cb_stats(user_id);