import time
from typing import List, Optional
import requests
from bs4 import BeautifulSoup, Tag
from urllib.parse import urljoin, urlparse
import logging
from readiness_probe import ReadinessProbe


class BS4Spider:
    def __init__(self, base_url: str) -> None:
        self.visited: set[str] = set()
        self.collected_data: List[str] = []  # Change type to List[str]
        self.base_domain: str = urlparse(base_url).netloc
        self.logger = logging.getLogger(__name__)

    def process_page(self, url: str) -> str:  # Change return type to str
        start_time: float = time.time()

        # Notifying the readiness probe that the service is alive
        ReadinessProbe().update_last_seen()

        # Fetch and parse the URL
        soup = self.fetch_and_parse(url)
        if not soup:
            return ""

        # Extract data from the page
        page_content = self.extract_data(soup)

        end_time: float = time.time()  # Record the end time
        elapsed_time: float = end_time - start_time
        self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

        return page_content  # Return as a single string

    def extract_links(self, start_url: str) -> List[str]:
        start_time: float = time.time()
        links_to_visit: List[str] = [start_url]
        all_links: List[str] = []

        while links_to_visit:
            url = links_to_visit.pop()
            if url in self.visited:
                continue

            # Notifying the readiness probe that the service is alive
            ReadinessProbe().update_last_seen()

            # Fetch and parse the URL
            soup: Optional[BeautifulSoup] = self.fetch_and_parse(url)
            if not soup:
                continue

            self.visited.add(url)
            all_links.append(url)

            # Extract all links from the page
            links: List[str] = [a['href'] for a in soup.find_all('a', href=True)]
            for link in links:
                # Convert relative links to absolute links
                absolute_link: str = urljoin(url, link)
                parsed_link = urlparse(absolute_link)
                # Check if the link is an HTTP/HTTPS link, not visited yet, and does not contain a fragment
                if parsed_link.scheme in ['http',
                                          'https'] and absolute_link not in self.visited and not parsed_link.fragment:
                    # Ensure the link is within the same domain
                    if parsed_link.netloc == self.base_domain:
                        links_to_visit.append(absolute_link)

        end_time: float = time.time()  # Record the end time
        elapsed_time: float = end_time - start_time
        self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

        return all_links

    def fetch_and_parse(self, url: str) -> Optional[BeautifulSoup]:
        try:
            self.logger.info(f"Processing URL: {url}")
            response = requests.get(url)
            if response.status_code == 200:
                soup: BeautifulSoup = BeautifulSoup(response.text, 'html.parser')
                return soup
            else:
                self.logger.error(f"❌ failed to retrieve URL: {url}, Status Code: {response.status_code}")
                return None
        except Exception as e:
            self.logger.error(f"❌ error fetching URL: {url}, Error: {e}")
            return None

    def extract_data(self, soup: BeautifulSoup) -> Optional[str]:
        elements: List[Tag] = soup.find_all(['p', 'article', 'div'])
        paragraphs: List[str] = []

        for element in elements:
            text: str = element.get_text(strip=True)

            if text and text not in paragraphs and len(text) > 10:
                paragraphs.append(text)

        formatted_text: str = '\n\n '.join(paragraphs)
        return formatted_text
