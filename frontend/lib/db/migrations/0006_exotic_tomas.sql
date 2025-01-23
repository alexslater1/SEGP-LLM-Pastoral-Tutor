ALTER TABLE "Message" ALTER COLUMN "content" SET NOT NULL;--> statement-breakpoint
ALTER TABLE "Message" ADD COLUMN "annotations" json;