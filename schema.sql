create table "matches"
(
    id          integer PRIMARY KEY,
    ip          VARCHAR(15) NOT NULL,
    started_at  timestamp   not null,
    map         VARCHAR(50) NOT NULL,
    rounds      smallint    NOT NULL,
    duration    integer     NOT NULL,
    won         bool        NOT NULL default false,
    inserted_at timestamp   not null default current_timestamp,
    UNIQUE (ip, started_at, map)
)

create table "users"
(
    id               bigint PRIMARY KEY,
    kills            integer   NOT NULL default 0,
    deaths           integer   NOT NULL default 0,
    kd               numeric(10, 2),
    all_weapon_stats jsonb     NOT NULL default ''{}''::jsonb,
    inserted_at      timestamp not null default current_timestamp
)

create table "user_stats"
(
    match_id     integer NOT NULL,
    user_id      bigint  NOT NULL,
    kills        integer NOT NULL default 0,
    deaths       integer NOT NULL default 0,
    weapon_stats jsonb   NOT NULL default ''{}''::jsonb,

    FOREIGN KEY (match_id) REFERENCES matches (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
)


