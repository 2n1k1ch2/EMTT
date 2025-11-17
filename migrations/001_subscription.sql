CREATE TABLE subscription (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    service_name VARCHAR(100) NOT NULL,
    price INT,
    start_date TIMESTAMP NOT NULL,
    finish_date TIMESTAMP,
    CONSTRAINT check_dates CHECK (finish_date IS NULL OR finish_date > start_date)
);
CREATE INDEX idx_subscription_user_id ON subscription(user_id);
CREATE INDEX idx_subscription_dates ON subscription(start_date, finish_date);