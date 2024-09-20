BEGIN;
    ALTER TABLE rooms ADD CONSTRAINT unique_room_constraint UNIQUE (room_name, room_key);

    ALTER TABLE events DROP COLUMN event_date;

    ALTER TABLE events ADD COLUMN start_date timestamp not null;
    ALTER TABLE events ADD COLUMN end_date timestamp not null;

COMMIT;