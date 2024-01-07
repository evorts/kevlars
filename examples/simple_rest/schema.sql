create table if not exists todo(
    id int primary key generated always as identity,
    title varchar(100),
    description varchar(1000),
    created_at timestamp default current_timestamp,
    updated_at timestamp
);