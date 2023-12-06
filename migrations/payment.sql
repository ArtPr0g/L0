create table payment
(
    payment_uuid uuid,
    transaction varchar(255),
    request_id varchar(255),
    currency varchar(255),
    provider varchar(255),
    amount integer,
    payment_dt integer,
    bank varchar(255),
    delivery_cost integer,
    goods_total integer,
    custom_fee integer
);