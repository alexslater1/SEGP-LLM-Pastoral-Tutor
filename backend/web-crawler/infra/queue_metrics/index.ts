// src/index.ts

import {
  CloudWatchClient,
  PutMetricDataCommand,
} from "@aws-sdk/client-cloudwatch";
import { createClient } from "redis";

// Lambda handler
export const handler = async (
  _event: unknown
): Promise<{ statusCode: number; body: string }> => {
  // Get Redis connection details from environment variables
  const redisHost = process.env.REDIS_HOST;
  const redisPort = process.env.REDIS_PORT || "6379";
  const redisPassword = process.env.REDIS_PASSWORD;
  const redisDatabase = process.env.REDIS_DB || "0";
  if (!redisHost) {
    console.error("Missing REDIS_HOST environment variable.");
    throw new Error("REDIS_HOST not defined");
  }

  console.log("Connecting to Redis", {
    redisHost,
    redisPort,
    redisPassword,
  });

  // Create a Redis client
  const client = createClient({
    socket: {
      host: redisHost,
      port: Number(redisPort),
      tls: true,
    },
    password: redisPassword,
    database: Number(redisDatabase),
  });

  client.on("error", (err) => {
    console.error("Redis Client Error", err);
  });

  try {
    // Connect to Redis
    await client.connect();

    // Query the lengths of the queues.
    // It is assumed that the keys for the queues are "urls" and "rag".
    console.log("Querying Redis for queue lengths");
    const [urlsQueueLength, ragQueueLength] = await Promise.all([
      client.lLen("urls"),
      client.lLen("rag"),
    ]);

    console.log(
      `Queue lengths: urls=${urlsQueueLength}, rag=${ragQueueLength}`
    );

    console.log("Creating CloudWatch client");
    // Create a CloudWatch client
    const cloudwatch = new CloudWatchClient({});

    console.log("Building metric data");
    // Build the metric data for both queues.
    const metricData = [
      {
        MetricName: "ScraperQueueLength",
        Dimensions: [
          {
            Name: "QueueName",
            Value: "urls",
          },
        ],
        Timestamp: new Date(),
        Value: urlsQueueLength,
        Unit: "Count" as const,
      },
      {
        MetricName: "RAGQueueLength",
        Dimensions: [
          {
            Name: "QueueName",
            Value: "rag",
          },
        ],
        Timestamp: new Date(),
        Value: ragQueueLength,
        Unit: "Count" as const,
      },
    ];

    console.log("Sending PutMetricData command");
    // Prepare and send the PutMetricData command
    const putMetricDataCommand = new PutMetricDataCommand({
      Namespace: "CustomRedisMetrics", // same as in your Terraform
      MetricData: metricData,
    });

    console.log("Sending PutMetricData command");
    await cloudwatch.send(putMetricDataCommand);
    console.log("Successfully published metrics to CloudWatch.");

    return {
      statusCode: 200,
      body: JSON.stringify({
        message: "Metrics updated",
        urlsQueueLength,
        ragQueueLength,
      }),
    };
  } catch (error) {
    console.error("Error in Lambda execution:", error);
    throw error;
  } finally {
    // Ensure the Redis client disconnects after processing
    await client.disconnect();
  }
};
