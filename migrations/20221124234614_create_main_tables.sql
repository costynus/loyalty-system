-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE public.order (
    id serial PRIMARY KEY,
    order_number text,
    status VARCHAR(32),
    accrual decimal,
    uploaded_at timestamp,
    user_id int,
    CONSTRAINT FK_order_user FOREIGN KEY (user_id) REFERENCES public.user (id)
);
CREATE TABLE balance (
    balance decimal,
    withdrawal decimal,
    user_id int,
    CONSTRAINT FK_balance_user FOREIGN KEY (user_id) REFERENCES public.user (id)
);
CREATE TABLE withdrawal (
    order_number text,
    sum_number decimal,
    updated_at timestamp,
    user_id int,
    CONSTRAINT FK_withdrawal_user FOREIGN KEY (user_id) REFERENCES public.user (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE public.order;
DROP TABLE public.balance;
DROP TABLE public.withdrawal;
-- +goose StatementEnd
