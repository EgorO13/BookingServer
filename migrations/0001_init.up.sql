CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       role TEXT NOT NULL CHECK (role IN ('admin', 'user')),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE rooms (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name TEXT NOT NULL,
                       description TEXT,
                       capacity INT,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE schedules (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
                           days_of_week INT[] NOT NULL,
                           start_time TIME NOT NULL,
                           end_time TIME NOT NULL,
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                           UNIQUE(room_id)
);

CREATE TABLE slots (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
                       start_time TIMESTAMP WITH TIME ZONE NOT NULL,
                       end_time TIMESTAMP WITH TIME ZONE NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_slots_room_start ON slots(room_id, start_time);

CREATE TABLE bookings (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          status TEXT NOT NULL CHECK (status IN ('active', 'cancelled')) DEFAULT 'active',
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_bookings_slot_status ON bookings(slot_id, status);
CREATE INDEX idx_bookings_user_status ON bookings(user_id, status);

INSERT INTO users (id, email, role) VALUES
                                        ('00000000-0000-0000-0000-000000000001', 'admin@mail.com', 'admin'),
                                        ('00000000-0000-0000-0000-000000000002', 'user@mail.com', 'user')
    ON CONFLICT (id) DO NOTHING;