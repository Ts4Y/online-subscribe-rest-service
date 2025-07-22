-- +goose Up
-- +goose StatementBegin
create table
   subscriptions (
      id uuid primary key,
      service_name text not null,
      price INT not null,
      user_id text not null,
      start_date date not null,
      end_date date
   );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE subscriptions;

-- +goose StatementEnd