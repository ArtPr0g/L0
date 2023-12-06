create table orders
(
    order_uid varchar(255),
    track_number varchar(255),
    entry varchar(255),
    delivery uuid,
    payment uuid,
    items uuid[],
    locale varchar(255),
    internal_signature varchar(255),
    customer_id varchar(255),
    delivery_service varchar(255),
    shardkey varchar(255),
    sm_id integer,
    date_created timestamp with time zone,
    oof_shard varchar(255)
);