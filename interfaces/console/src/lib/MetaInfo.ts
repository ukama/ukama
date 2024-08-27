import { IP_API_BASE_URL, IPFY_URL } from '@/constants';

const getMetaInfo = async () => {
  return await fetch(IPFY_URL, {
    method: 'GET',
  })
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .then((data) =>
      fetch(`${IP_API_BASE_URL}/${data.ip}/json/`, {
        method: 'GET',
      }),
    )
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .catch((err) => {
      console.log(err);
      return {};
    });
};

export { getMetaInfo };
