import requests
import time
from concurrent.futures import ThreadPoolExecutor


def send_request(url, data):
    try:
        headers = {'Content-Type': 'application/json'}
        response = requests.post(
            url=url,
            json=data,
            headers=headers,
        )
        if response.status_code == 201:
            return 0
        else:
            return response.json()
    except requests.exceptions.RequestException as e:
        return str(e)


def main():
    url = "http://localhost:9000/"
    number_of_requests = 100
    data = {
        "expireAt": "2025-02-08T09:20:41Z",
        "url": "https://google.com"
    }

    start_time = time.time()

    with ThreadPoolExecutor(max_workers=10) as executor:
        responses = list(executor.map(
            lambda x: send_request(url, data), range(number_of_requests))
        )

    print(responses)

    end_time = time.time()
    duration = end_time - start_time

    success_count = responses.count(0)
    failure_count = number_of_requests - success_count

    print(f"Total requests: {number_of_requests}")
    print(f"Successful requests: {success_count}")
    print(f"Failed requests: {failure_count}")
    print(f"Total time taken: {duration} seconds")


if __name__ == "__main__":
    main()
