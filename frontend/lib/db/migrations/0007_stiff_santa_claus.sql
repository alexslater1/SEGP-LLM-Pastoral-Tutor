ALTER TABLE "User" ADD COLUMN "role" text DEFAULT 'student' NOT NULL;--> statement-breakpoint
ALTER TABLE "User" ADD COLUMN "created_at" timestamp DEFAULT now() NOT NULL;