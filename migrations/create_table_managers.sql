
create table if not customers
(
    id bigserial primary key,
    name	text not null,
    phone 	text 	not null unique,
    password text 	not null,
    active 	boolean not null default true,
    created TIMESTAMP not null default current_timestamp 
);

create table if not managers 
(
    id bigserial primary key,
    name	text not null,
    salary integer not null default 0,
    plan    integer not null default 0,
    boss_id bigint references managers,
    departament text,
    phone 	text 	not null unique,
    password text 	,
    is_admin boolean not null default true,
    active 	boolean not null default true,
    created TIMESTAMP not null default current_timestamp 
);

create table if not customers_tokens 
(
    token text not null unique,
    customer_id bigint not null references customers,
    expire  TIMESTAMP not null default current_timestamp + interval '1 hour',
    created TIMESTAMP not null default current_timestamp
);

create table if not managers_tokens 
(
    token text not null unique,
    manager_id bigint not null references managers,
    expire  TIMESTAMP not null default current_timestamp + interval '1 hour',
    created TIMESTAMP not null default current_timestamp
);

create table if not products 
(
    id      bigserial primary key,
    name    text not null,
    price   integer not null CHECK (price >0),
    qty     integer not null default 0 CHECK (qty >=0),
    active 	boolean not null default true,
    created TIMESTAMP not null default current_timestamp 
);

create table if not sales 
(
    id          bigserial primary key,
    manager_id  bigint not null references managers,
    customer_id bigint not null,
    created     TIMESTAMP not null default current_timestamp 
);

create table if not sales_positions 
(
    id          bigserial primary key,
    product_id  bigint not null references products,
    sale_id  bigint not null references sales,
    price integer not null CHECK (price >= 0),
    qty     integer not null default 0 CHECK (qty >=0),
    created     TIMESTAMP not null default current_timestamp 
);