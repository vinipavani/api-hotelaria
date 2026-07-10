CREATE TABLE IF NOT EXISTS rooms (
    id BIGSERIAL PRIMARY KEY,
    hotel_id BIGINT NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    number VARCHAR(4) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('single', 'double', 'suite')),
    capacity INT NOT NULL CHECK (capacity > 0),
    per_night_value NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_hotel_room_number UNIQUE (hotel_id, number)
);