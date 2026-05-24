import { Body, Container, Head, Heading, Html, Tailwind, Text, Img, Row, Column, Hr, Link } from 'react-email'

export const NotionMagicLinkEmail = () => (
  <Html>
    <Head />
    <Tailwind>
      <Body className="bg-white">
        <Container className="mx-auto px-3">
          <Heading className="my-10 p-0 font-sans font-bold text-[#333333] text-xl">EMAIL-TITLE</Heading>
          <Text>EMAIL-INTRO</Text>
          <Row className='rounded-xl bg-gray-200 p-3'>
            <Column className='rounded-xl bg-gray-100' align="center">
              <Img src="https://github.com/microsoft/fluentui-emoji/blob/main/assets/Spiral%20notepad/3D/EMAIL-FILEICON_3d.png?raw=true" className='size-10' style={{ display: 'block', margin: 'auto' }} />
            </Column>
            <Column className='w-full pl-3'>
              <Text className='m-0'>
                EMAIL-FILENAME
              </Text>
            </Column>
          </Row>
          <Row className='my-10 text-center'>
            <Column className='w-1/2'>
              <Text className='m-0 text-xl font-bold'>
                EMAIL-IP
              </Text>
              <Text className='text-sm opacity-50 m-0'>
                IP
              </Text>
            </Column>
            <Column className='w-1/2'>
              <Text className='m-0 text-xl font-bold'>
                EMAIL-REGION
              </Text>
              <Text className='text-sm opacity-50 m-0'>
                EMAIL-LABEL-REGION
              </Text>
            </Column>
          </Row>
          <Hr className="my-4 border-gray-300 border-t-2" />
          <Text className='m-0 text-sm'>
            Powered by <Link href="https://github.com/keven1024/015" className='text-blue-500'>015</Link>
          </Text>
        </Container>
      </Body>
    </Tailwind>
  </Html>
)

export default NotionMagicLinkEmail
