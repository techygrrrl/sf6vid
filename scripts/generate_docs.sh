#!/usr/bin/env sh

go install

rm USAGE.md

echo "# Usage" >> USAGE.md
echo "## sf6vid" >> USAGE.md
echo '```' >> USAGE.md
sf6vid --help >> USAGE.md
echo '```' >> USAGE.md

echo "### censor" >> USAGE.md
echo '```' >> USAGE.md
sf6vid censor --help >> USAGE.md
echo '```' >> USAGE.md

echo "### trim" >> USAGE.md
echo '```' >> USAGE.md
sf6vid trim --help >> USAGE.md
echo '```' >> USAGE.md

echo "### shrink" >> USAGE.md
echo '```' >> USAGE.md
sf6vid shrink --help >> USAGE.md
echo '```' >> USAGE.md
