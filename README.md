# Finetic

Finetic is a personal finance app designed to turn SMS-based financial data into useful spending insights.

The project is being built in stages, starting with app setup, transaction parsing, and local development infrastructure.

## Goals

- Track income and expenses
- Parse SMS transaction messages into structured data
- Categorize spending patterns
- Surface useful financial insights over time
- Build a clean foundation for future product features

## Tech stack

- Flutter / Dart for the mobile app
- Go for backend services
- Python for AI-related processing and forecasting
- PostgreSQL for data storage
- Docker for local environment setup

## Repository structure

```text
finetic/
├── apps/
│   └── lib/                 # Flutter app
├── backend/                 # Go backend
├── ai-service/              # Python AI service
├── docs/                    # Project documentation and planning
├── scripts/                 # Setup and utility scripts
├── docker-compose.yml       # Local service orchestration
├── .gitignore               # Ignore generated and local files
├── LICENSE                  # Project license
├── README.md                # Project overview
└── .github/                 # CI/CD workflows
```

## Getting started

### Prerequisites

Make sure these are installed:

- Flutter SDK
- Dart SDK
- Go
- Python
- Docker
- Git

### Install app dependencies

```bash
cd apps/lib
flutter pub get
```

### Run the app

```bash
cd apps/lib
flutter run
```

## Current status

This project is in active development. The app foundation and tooling are being established before expanding into more advanced features.

## Planned features

- SMS parsing and transaction extraction
- transaction storage and retrieval
- categorization and reporting
- dashboard and analytics views
- AI-based spending insights and forecasting
- deployment and production setup

## Contributing

This project is still evolving, so contributions are expected to stay focused and incremental.

Typical workflow:

```bash
git checkout -b feature/your-feature
git add .
git commit -m "Add your feature"
git push origin feature/your-feature
```

## Notes

This README is intentionally simple and will be expanded as the product and infrastructure mature.
