CREATE TABLE "OnboardingSession" (
    "session_id" uuid NOT NULL,
    "scenario_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "status" text NOT NULL,
    "started_at" timestamp NOT NULL,
    "finished_at" timestamp,
    PRIMARY KEY ("session_id")
);
