CREATE TABLE "OnboardingEvent" (
    "event_id" uuid NOT NULL,
    "session_id" uuid NOT NULL,
    "step_id" uuid,
    "type" text NOT NULL,
    "data" jsonb NOT NULL,
    "occurred_at" timestamp NOT NULL,
    "received_at" timestamp NOT NULL,
    PRIMARY KEY ("event_id"),

    CONSTRAINT "fk_OnboardingEvent_session_id"
      FOREIGN KEY ("session_id")
      REFERENCES "public"."OnboardingSession" ("session_id")
);
CREATE INDEX "OnboardingEvent_index_1" ON "OnboardingEvent" ("session_id");
