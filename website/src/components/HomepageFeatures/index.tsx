import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  icon: string;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Multi-Protocol Support',
    icon: '🔐',
    description: (
      <>
        Test all modern DNS protocols in one tool: <strong>Do53</strong> (UDP/TCP),
        <strong> DoT</strong>, <strong>DoH</strong>, and <strong>DoQ</strong>.
      </>
    ),
  },
  {
    title: 'High Performance',
    icon: '⚡',
    description: (
      <>
        Built with <strong>Go</strong> for speed. Query multiple DNS servers
        concurrently with minimal latency.
      </>
    ),
  },
  {
    title: 'Production Ready',
    icon: '📊',
    description: (
      <>
        Built-in <strong>Prometheus</strong> metrics, rate limiting, and
        async task processing with Redis.
      </>
    ),
  },
  {
    title: 'REST API & CLI',
    icon: '🔌',
    description: (
      <>
        Use the <strong>CLI</strong> for quick tests or the <strong>REST API</strong>
        for integration and automation.
      </>
    ),
  },
  {
    title: 'Easy Deployment',
    icon: '🐳',
    description: (
      <>
        Deploy with <strong>Docker Compose</strong> in seconds.
        Horizontal scaling with worker pools.
      </>
    ),
  },
  {
    title: 'Rate Limiting',
    icon: '🛡️',
    description: (
      <>
        Protect your infrastructure with configurable <strong>per-IP
        rate limiting</strong> and burst control.
      </>
    ),
  },
];

function Feature({title, icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <div className={styles.featureIcon}>{icon}</div>
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
