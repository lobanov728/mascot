CREATE TABLE IF NOT EXISTS transaction(
   transaction_uuid uuid PRIMARY KEY,
   transaction_ref TEXT NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
   caller_id INT NOT NULL,
   player_uuid TEXT NOT NULL,
   withdraw INT NOT NULL,
   deposit INT NOT NULL,
   balance INT NOT NULL,
   status VARCHAR(64) NOT NULL,
   currency VARCHAR (8) NOT NULL,
   game_id VARCHAR (255),
   game_round_ref VARCHAR (255),
   source TEXT,
   reason TEXT,
   session_id VARCHAR (255),
   session_alternative_id VARCHAR (255),
   spin_details JSONB,
   bonus_id VARCHAR (255),
   charge_free_rounds INT,
   free_rounds_left INT
);

CREATE INDEX transaction_clusterd_idx
  ON transaction (caller_id, player_uuid, transaction_ref);