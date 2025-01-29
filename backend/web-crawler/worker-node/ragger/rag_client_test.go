package ragger

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactsFor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Contact
	}{
		{
			name:  "should extract email",
			input: "You can reach me at john.doe@example.com for more information.",
			expected: []Contact{
				{
					Value: "john.doe@example.com",
					Type:  "email",
				},
			},
		},
		{
			name:  "should extract phone number",
			input: "Call me at +1 (555) 123-4567 anytime.",
			expected: []Contact{
				{
					Value: "+1 (555) 123-4567",
					Type:  "phone",
				},
			},
		},
		{
			name:  "should extract URL",
			input: "Visit our website at https://example.com/contact for details.",
			expected: []Contact{
				{
					Value: "https://example.com/contact",
					Type:  "url",
				},
			},
		},
	}

	client := &RagClient{} // No need to initialize with real model for contact extraction

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ContactsFrom(tt.input)
			if err != nil {
				t.Errorf("ContactsFor() error = %v", err)
				return
			}

			if len(got) != len(tt.expected) {
				t.Errorf("ContactsFor() got %d contacts, expected %d", len(got), len(tt.expected))
				return
			}

			for i, contact := range got {
				if contact.Value != tt.expected[i].Value {
					t.Errorf("ContactsFor() got value = %v, expected %v", contact.Value, tt.expected[i].Value)
				}
				if contact.Type != tt.expected[i].Type {
					t.Errorf("ContactsFor() got type = %v, expected %v", contact.Type, tt.expected[i].Type)
				}
			}
		})
	}
}

func TestContactsFor2(t *testing.T) {
	client := NewRagClient("../model", "../libonnxruntime.so.1.20.1")

	text := `hello world`

	contacts, err := client.ContactsFrom(text)
	if err != nil {
		t.Fatalf("error getting contacts %v", err)
	}

	log.Printf("contacts: %+v", contacts)
}

func TestChunksFor(t *testing.T) {
	client := NewRagClient("../model", "../libonnxruntime.so.1.20.1")

	text := `# Ethan Hosier

Third Year Computing Student at Imperial College London, passionate about writing software.

[London, UK](https://www.google.com/maps/search/london/@51.4850816,-0.1998848,14z?entry=ttu)

[LinkedIn](https://www.linkedin.com/in/ethan-hosier-623474253/) [CV](/resume.pdf)

[ethanjhosier@gmail.com](mailto:ethanjhosier@gmail.com)

EH

## About

Looking for job opportunities, particularly in early stage startups

## Work Experience

### Revolut Internship

July 2024-Sept 2024

#### Software Engineer (Intern)

Designed and implemented the backend of an Authorised Data Source that consolidates product and pricing rule configurations into a single interface, eliminating the need for manual database queries—used in production daily by the CEO and leadership team. Leveraged concurrency techniques and caching to reduce query latency by 99%, significantly improving system performance and efficiency. Integrated with multiple internal services, including the Terms Service to assign legal documents to tax configurations, and the Revolut People Platform to incorporate employee details into the admin app, ensuring seamless interaction across different systems

Technologies:
Java

TypeScript

React

### Mia Marketing Founder

2024

#### Co-Founder, CTO

Led the development of an MVP for a fully autonomous AI-driven digital marketing platform, attracting interest from over 40 different companies, including those with annual revenues exceeding £1 million. Built a scalable, microservices-based infrastructure on AWS and Azure of proxies, web crawlers and scrapers. Developed automated workflows for data-driven social media content creation, posting, and lead generation through personalized email outreach campaigns. Awarded £2,500 in prize money from Imperial College London’s Entrepreneurship Journey, alongside £3,500 in Azure credits from the Microsoft Startup Hub

Technologies:
Go

Python

Typescript

React.js

AWS

Microsoft Azure

Docker

### 321 Restaurants Founder

2023

#### Full Stack Developer / Founder

Founded a software startup dedicated to empowering small, independent restaurants through affordable solutions. Conducted market research to understand industry needs and, as the sole full stack developer, developed a comprehensive content management system, dashboard and individualized websites. Gained valuable experience in product development, market analysis, and demonstrated proficiency in designing and implementing full stack solutions.

Technologies:
TypeScript

Node.js

React/Next.js

Firebase

## Education

### Imperial College London

2022 \- 2025

Bachelor's Degree in Computing

## Skills

Java

Python

C

Go

Haskell

Scala

TypeScript

JavaScript

AWS

Microsoft Azure

Docker

SQL

React

Next.js

React Native

Node.js

Pandas

Pytorch

## Hackathons

### [IC Hack](https://github.com/fm0ss/financial_tool_ichack)  24H

February 2024

#### AI Trading Statements Analyser

As part of a team of five Imperial Computing students, developed an AI powered tool that performed web scraping and analysis of trading statements from the London Stock Exchange website. Our tool successfully navigated the statements' biased and overly positive corporate language, providing fact checking and context to the user by searching the web with LangChain. As an individual, I trained a custom BERT language model to perform sentiment analysis on the positivity of each sentence within the statement. I also built out the entire frontend of the app, as well as some of the backend.

Technologies:
Python

Flask

TypeScript

React/Next.js

LangChain

Docker

Hugging Face

OpenAI

### [AI Ventures Hackathon](https://github.com/EthanHosier/ai-hack/)  24H

January 2024

#### GILO Consulting

Worked with two business students to create an AI-powered consulting platform, which we pitched to VCs. Our platform aimed to tackle issues faced such as the high costs of traditional consulting, limited access to expertise and support, difficulty finding the right experts, and a lack of transparency and accountability. As the sole technical member of our team, developed gainshare contract generation, consultant profiles, a dashboard, and a chatbot custom-trained to match our clients directly to the most relevant consultants. Gained valuable experience working in a team with diverse skillsets and producing software quickly.

Technologies:
TypeScript

React/Next.js

OpenAI

## Projects

### Wacc Compiler

Compiler for the WACC programming language that translates WACC source code into x86-64 machine code, handling lexical analysis, semantic analysis, code generation, and optimization. Web-based IDE built with Monaco editor and Next.js, supporting WACC code writing, compilation, and execution with syntax and error highlighting.

University Project

Group project

Scala

Typescript

Next.js

### Groupify

Mobile app for university students to find project groups, inspired by Tinder, featuring a custom matching algorithm, project creation, recommendations, and integrated chat rooms. Collaborated with stakeholders to define and refine requirements, ensuring user needs alignment.

University Project

Group project

Typescript

React Native

### [Large Language Model (GPT)](https://github.com/EthanHosier/LLM)

github.comEthanHosier/LLM

Large language model from scratch, trained on the OpenWebText dataset. Capable of generating human-like text using next token prediction.

Side Project

Python

Pytorch

Jupyter Notebook

### [Restaurant Dashboard & CMS](https://github.com/EthanHosier/Restaurant-Dashboard-CMS)

github.comEthanHosier/Restaurant-Dashboard-CMS

Full stack dashboard enabling customer profiling, bookings and venue management. Includes a content management system for populating websites. Open-sourced from my startup.

Startup

TypeScript

React/Next.js

Firebase

### [Pop Up!](https://github.com/EthanHosier/smDemo)

github.comEthanHosier/smDemo

MVP for my idea of a proximity-based social media mobile app. Profiles of nearby app users dynamically pop up on users' phones as they walk past each other, facilitating real-time interactions.

Side Project

MVP

JavaScript

React Native

Bluetooth

Firebase

### [ARMv8 AArch64 – Assembler and Emulator](https://github.com/EthanHosier/armv8-compilor-emulater)

github.comEthanHosier/armv8-compilor-emulater

AArch64 emulator that simulates the execution of an AArch64 binary file on a Raspberry Pi, and AArch64 assembler that translates an AArch64 assembly source file containing A64 instructions into a binary file that can subsequently be executed by the emulator.

University Project

Group Project

C

### [RSS Aggregator](https://github.com/EthanHosier/rssagg)

github.comEthanHosier/rssagg

RSS feed web scraper with a REST API and authentication using api keys and middleware. Stores scraped results in a postreSQL database.

Side Project

Golang

PostgreSQL

REST

### [Trello Clone](https://github.com/EthanHosier/Trello-clone)

github.comEthanHosier/Trello-clone

Full stack Trello clone with authentication, organizations, workspaces, stripe subscriptions and CRUD operations implemented with server actions.

Side Project

TypeScript

Next.js

Prisma

MySQL

### [Digit Recognition Neural Network (No ML/Maths Libraries)](https://github.com/EthanHosier/mnist-digits-nn)

github.comEthanHosier/mnist-digits-nn

Deep multilayer perceptron neural network for handwritten digit recognition, without any machine learning or maths libraries. Implemented backpropagation from scratch.

Side Project

Java

Machine Learning

### [Airbnb Clone](https://github.com/EthanHosier/airbnb)

github.comEthanHosier/airbnb

Airbnb mobile app clone, developed with React Native.

Side Project

Typescript

React Native

Expo

Press "⌘J" to open the command menu`

	chunks, err := client.ChunksFrom(text)
	if err != nil {
		t.Fatalf("error getting chunks %v", err)
	}

	assert.Equal(t, 6, len(chunks))

	for _, chunk := range chunks {
		log.Printf("chunk: %v \n=============\n", chunk)
	}
}

func TestEmbeddingsFor(t *testing.T) {
	client := NewRagClient("../model", "../libonnxruntime.so.1.20.1")

	embeddings, err := client.EmbeddingsFor("hello world")
	if err != nil {
		t.Fatalf("error getting embeddings %v", err)
	}

	assert.Equal(t, 384, len(embeddings))
	log.Printf("embeddings (first 10): %v", embeddings[:10])
}
