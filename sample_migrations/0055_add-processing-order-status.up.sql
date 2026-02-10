DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_enum
        WHERE enumlabel = 'processing'
        AND enumtypid = 'order_status'::regtype
    ) THEN
        ALTER TYPE order_status ADD VALUE 'processing' BEFORE 'confirmed';
    END IF;
END $$;
