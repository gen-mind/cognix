# # ###################################
# # Distroless working
# # ###################################
# FROM python:3.11.7-slim  AS build-env

# WORKDIR /app
# ADD ./ ./ 

# # Install any needed packages specified in requirements.txt
# RUN pip install --upgrade pip && \
#     pip install --no-cache-dir -r requirements.txt \
#     && rm -rf /root/.cache

# # Install dependencies and clean up
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     && rm -rf /var/lib/apt/lists/*

# # Use a distroless Python image as a parent image
# FROM gcr.io/distroless/python3-debian12

# # Copy the application code and installed packages from the build environment
# COPY --from=build-env /app /app
# COPY --from=build-env /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages

# # Set the working directory in the container
# WORKDIR /app

# # Set environment variable for Python path
# ENV PYTHONPATH=/usr/local/lib/python3.11/site-package

# # Run the application
# CMD ["semantic_service.py"]





# ###################################
# python slim working
# ###################################
# using Python 3.11.0 as the newr supported by PyTorch
FROM python:3.11.7-slim-bookworm


# Set the working directory inside the container
WORKDIR /app
ADD ./ ./

# Copy the requirements.txt first to leverage Docker cache
COPY semantic/requirements.txt /app/

# Install dependencies and clean up
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*


# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir -r requirements.txt \
    && rm -rf /root/.cache

# the container before this 3 lines is 107mb, after it's 982mb!!!!
#RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
#RUN pip install transformers==4.40.2 tqdm==4.66.4 numpy==1.26.4 scikit-learn==1.4.2 scipy==1.13.0 nltk==3.8.1 pillow==10.3.0

# Copy the rest of the application code
COPY semantic /app/semantic
COPY lib_py/cognix_lib /app/cognix_lib

# Set the PYTHONPATH environment variable
ENV PYTHONPATH="/app"

# Command to run your application
CMD ["python", "semantic/semantic_service.py"]













# # Use a slim Python base image
# # using Python 3.11.0 as the newr supported by PyTorch
# FROM python:3.11.7-slim-bookworm


# # Set the working directory inside the container
# WORKDIR /app

# # Copy the requirements.txt first to leverage Docker cache
# COPY requirements.txt .

# # Install dependencies and clean up
# RUN apt-get update && apt-get install -y --no-install-recommends \
#     build-essential \
#     && rm -rf /var/lib/apt/lists/*
# RUN pip install --upgrade pip
# RUN pip install --no-cache-dir torch==2.3.0 --index-url https://download.pytorch.org/whl/cpu
# RUN pip install --no-cache-dir -r requirements.txt \
#     && rm -rf /root/.cache

# # Copy the rest of the embedder application code
# COPY . .

# # Command to run your application
# CMD ["python", "semantic.py"]
