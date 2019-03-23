from setuptools import setup, find_packages

setup(
    name='go-webhdfs-run',
    version='1.0.0',
    packages=find_packages(),
    install_requires=[
        'docker',
    ]
)